# tryco-backend
> TryCo Backend API

코인 모의 투자 서비스 TryCo의 백엔드 리파지토리입니다.

## 프로젝트 셋업

언어는 Go, 프레임워크는 Fiber를 사용합니다.

로컬에서 API 호출은 vscode의 REST Client extention으로 테스트합니다. (`/test.http` 파일 이용)

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