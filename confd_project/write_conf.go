// 用golang编写，可以应付多种环境，并且可以对文件编译，防止密钥泄露
package main

import (
 "encoding/json"
 "flag"
 "fmt"
 "io/ioutil"
 "os"

 "github.com/kirinlabs/HttpRequest"
)

const (
 authKey = "eyJ0eXAiOiJKV1QiLCJ2MCIsImRhdGEiOnsidXNlcl9pZCI6MjgsInVzZXJuYW1lIjoic3MtdGVzdhbGciOiJIUzI1NiJ9.eyJleHAiOjE2NTQ4MjzdXBlcnVzZXIiOmZhbHNlfc2MDYsIm5iZiI6MTU1OTc4NzU4NiwiaWF0IjoxNTU5Nzg3NTk2LCJpc3MiOiJhdXRoOiBzcyIsInN1YiI6Im15IHRva2VuIiwiaWQiOiIxNTYxODcxODACIsIm5pY2tuYWJpc19X0.gMGMRKqtd_CM6rIzE8mxuwR8c8dz_hyH21FETOO4XbE"
 confURL = "https://codo.domain.com/api/kerrigan/v1/conf/publish/config/"
)

var (
 project_code string
 environment  string
 service      string
 filename     string
 realfilename string
 configData   string
)

func getArgs() error {
 flag.StringVar(&project_code, "p", "ss", "项目代号")
 flag.StringVar(&environment, "e", "dev", "环境")
 flag.StringVar(&service, "s", "app", "应用名称")
 flag.StringVar(&filename, "f", "settings.py", "文件名")
 flag.StringVar(&realfilename, "r", "/tmp/settings.py", "最终写入文件")
 flag.Parse()
 if project_code == "ss" {
  fmt.Println("[Error] 参数不正确，请使用参数 --help 查看帮助")
  os.Exit(-5)
 }
 return nil
}

func writeWithIoutil(name, content string) {
 data := []byte(content)
 if ioutil.WriteFile(name, data, 0644) == nil {
  fmt.Println("[Success] 修改配置成功")
 }
}

func main() {
 getArgs()
 req := HttpRequest.NewRequest()

 // 设置超时时间，不设置时，默认30s
 req.SetTimeout(30)

 // 设置Headers
 req.SetHeaders(map[string]string{
  "Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
 })

 // 设置Cookies
 req.SetCookies(map[string]string{
  "auth_key": authKey,
 })

 // GET
 resp, err := req.Get(confURL, map[string]interface{}{
  "project_code": project_code,
  "environment":  environment,
  "service":      service,
  "filename":     filename,
 })

 if err != nil {
  fmt.Println("[Error]", err)
  // log.Println(err)
  os.Exit(-1)
 }

 if resp.StatusCode() == 200 {
  body, err := resp.Body()

  if err != nil {
   fmt.Println("[Error]", err)
   os.Exit(-2)
  }

  res := make(map[string]interface{})
  json.Unmarshal(body, &res)

  for _, v := range res["data"].(map[string]interface{}) {
   configData = fmt.Sprintf("%v", v)
  }
 } else {
  os.Exit(-3)
 }
 writeWithIoutil(realfilename, configData)
}

// 执行: go run write_conf.go -p code-v1 -e dev -s codo-admin -f settings.py -r /tmp/settings.py1 
