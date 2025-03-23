# tryco-backend
> TryCo Backend API

코인 모의 투자 서비스 TryCo의 백엔드 리파지토리입니다.

## 적용된 기술 스택 설명

언어는 Go, 프레임워크는 Fiber를 사용합니다.

로컬에서 API 호출은 vscode의 REST Client extention으로 테스트합니다. (`/test.http` 파일 이용)

혹은, 생성된 Swagger 문서 상에서 테스트해보는 것도 좋습니다.

### Swagger 도입

추후 프론트엔드와의 협업을 위해 Swagger 문서를 도입하였습니다.

Go에서 Swagger를 만들기 위해서는 `swag` 명령어를 설치한 다음 사용합니다. 또한 주석을 바탕으로 자동으로 swagger 문서를 만들어주는 미들웨어를 사용하였습니다. 따라서 먼저 아래와 같이 필수 주석을 간단하게 작성해줍니다.

```go
// @title 스웨거 타이틀
// @version 1.0
func main() {
    //...
}
```

이후 `gofiber/swagger` 미들 웨어를 설치하고 기본 핸들러를 등록해줍니다.

```go
func main() {
    //...
    app.Get("/swagger/*", swagger.HandlerDefault)
    //...
}
```

이제 `swag` 명령어를 설치하고, `swag fmt`와 `swag init`을 통해 문서를 생성합니다.

```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag fmt
swag init
```

> ⚠ 참고로 여기서 `swag` 명령어를 못찾는다는 에러가 나오면, macOS에서 `homebrew`로 `go`를 설치했을 경우 그럴 수 있다. 이럴 때는 `PATH` 환경 변수에 `~/go/bin` 경로를 추가해주면 된다. 혹은 직접 `~/go/bin/swag`을 실행해줘도 된다.

자동으로 생성된 문서를 `import` 해줍니다.

```go
import (
    //...
    _ "github.com/my-github-account/repo-name/docs" // 생성된 경로에 맞게 수정
)
```

> 또한 만약 API Operations쪽에 주석을 분명히 썼는데 목록에 나오지 않는다면, @router 주석을 썼는지 확인해볼 것. 필수임.

### Sqlc를 ORM 대신 사용

Sqlc는 Sql을 통해 코드를 생성해주는 도구입니다. 스키마와 쿼리문을 직접 `.sql`로 짜놓으면 그에 맞춰서 코드를 생성해줍니다.

가장 먼저 프로젝트 루트에 `sqlc.yaml` 파일을 만듭니다. (`.json`도 가능)

```yaml
version: "2"
sql:
  - engine: "postgresql"      # 사용할 DB 종류
    queries: "db/query.sql"   # 쿼리문 경로
    schema: "db/schema.sql"   # 스키마 경로(혹은 Migrations)
    gen:
      go:
        out: "db/sqlc"        # 코드 생성 결과물 경로
        sql_package: "pgx/v5" # 사용할 sql 패키지(database/sql 등)
        emit_json_tags: true  # 모델에 json 태그명을 붙일 것인지

```

그리고 위 설정에서 명시한 대로 `db/query.sql`과 `db/schema.sql`을 작성합니다. 이후에는 `sqlc generate` 명령을 통해 코드를 생성할 수 있습니다. 이때 query쪽에는 주석을 통해 생성할 코드의 이름 등을 지정할 수 있습니다. (`db/query.sql` 참고)

본 프로젝트에서는 Postgres를 사용하고, 드라이버로는 `pgx/v5`를 사용합니다. 따라서 아래와 같이 `main.go`에서 데이터베이스를 연결하고 생성된 `conn` 구조체를 `Queries` 구조체에 주입하여 사용할 수 있습니다.

```go
conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatalf("Unable to connect to database: %v\n", err)
}
defer conn.Close(context.Background())
queries := sqlc.New(conn)
```

이렇게 준비된 `queries` 구조체에는 `db/query.sql`에서 정의한 쿼리문이 함수로 들어가 있습니다. 이어서 이 구조체를 컨트롤러에서 사용할 수 있도록 핸들러(컨트롤러) 구조체를 만들고 주입합니다.

```go
// sqlc의 Queries를 주입 받는 핸들러(컨트롤러)
type Handler struct {
	queries *sqlc.Queries
}

func New(queries *sqlc.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}
```

이렇게하면 라우트 핸들러(컨트롤러)에서 쿼리문 함수를 사용할 수 있습니다. 더 자세한 예시는 `handlers/example.go`를 참고해주세요.

### 올바른 swagger model을 위한 sqlc type override

만약 위에서 설정한 대로 `sqlc generate`를 했다면 아래와 같이 모델이 생성된다.

```go
type Book struct {
	ID        int32            `json:"id"`
	Title     string           `json:"title"`
	Author    string           `json:"author"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
```

이때 날짜를 나타내는 `CreatedAt`의 타입이 `pgtype.Timestamp`인 것을 볼 수 있다. 이 때문에 swagger 문서에서 이 `Book` 모델의 모양을 보면, 실제 JSON 응답의 모양이 아닌 구조체의 모양이 나오게 된다.

이를 방지하기 위해서는 일반적인 Go의 내장 타입인 `time.Time`으로 정의 되도록 설정을 변경해야하는데, 아래와 같이 sqlc 설정을 바꿔주면 된다.

```yml
version: "2"
sql:
  # ... 생략
    gen:
      go:
        # ... 생략
        overrides:
          - db_type: "pg_catalog.timestamp"
            go_type: "time.Time"
            nullable: true
```

이렇게 하고 다시 `sqlc generate`를 해주면 `pgtype.Timestamp`가 아닌 `time.Time`을 쓰기 때문에, swagger 문서 상에서도 모델이 JSON 응답과 일치하게 된다.