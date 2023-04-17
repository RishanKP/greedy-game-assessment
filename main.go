package main

import(
  "greedy-games-assessment/api"
  "github.com/gin-gonic/gin"
  "strings"
  "os"
  "fmt"
)

type Command struct{
  Command string `json:"command"`
}

func handlerFunction(c *gin.Context){
  var cmd Command
  c.ShouldBindJSON(&cmd)

  strs := strings.Split(cmd.Command," ")
  fmt.Println(strs)
  err,response,value := api.ProcessCommand(strs)

  if err != nil{
    c.JSON(response,gin.H{
      "error":err.Error(),
    })
    return
  }

  if value != ""{
    c.JSON(response,gin.H{
      "value":value,
    })
    return
  } 

  c.JSON(response,gin.H{
    "message":"success",
  })
}

func main() {
  r := gin.New()

  r.POST("/",handlerFunction)
  
  port := os.Getenv("port")
  if port == ""{
    port = "8080"
  }

  r.Run(":" + port)
}

