# Goren
Generate Gore code from OAS3

## Usage
```
go install git.tenvine.cn/backend/goren/cmd/goren@latest
goren openapi.yaml
```

会自动生产以下文件：
```
├─goren
│  │  goren.go
│  │
│  └─api
│     api.go
│
└─main.go
```

随后执行：
```
go mod tidy
```

便可中`goren.go`和`main.go`中进行开发了。