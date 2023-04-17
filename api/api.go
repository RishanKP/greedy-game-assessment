package api
import (
  "strconv"
  "time"
  "fmt"
  "errors"
)

type mydb struct {
  value   string
  expires time.Time
  exists  bool
}

var myMap = make(map[string]mydb)
var myQueue = make(map[string][]int)

func HandleSetCommand(strs []string)(error,int,string){
  _value := strs[2]

  if len(strs) > 3{
    if strs[3] == "EX"{
      seconds,err := strconv.Atoi(strs[4])
      if err != nil{
        return errors.New("invalid command"),400,""
      }

      if len(strs) > 4 {
        err := checkExistence(strs[1],strs[5])
        if err != nil{
          return err,400,""
        }
      }
      myMap[strs[1]] = mydb{value: _value, expires: time.Now().Add(time.Second * time.Duration(seconds) ), exists: true}
    }else if strs[3] == "XX" || strs[3] == "NX"{
      err := checkExistence(strs[1],strs[5])
      if err != nil{          
        return err,400,""
      }
      myMap[strs[1]] = mydb{value: _value, expires: time.Now().AddDate(100,0,0), exists: true}
    }else{
      return errors.New("invalid command"),400,""
    }
  }else{
    myMap[strs[1]] = mydb{value: _value, expires: time.Now().AddDate(100, 0, 0), exists: true}
  }

  return nil,201,""
}

func Enqueue(strs []string)(error,int,string){
  for i := 1; i < len(strs) ; i++ {
    a,err := strconv.Atoi(strs[i])
    if err != nil{
      continue //appends only integer values
    }
    myQueue[strs[0]] = append(myQueue[strs[0]],a)  
  }
  fmt.Println(myQueue[strs[0]])
  
  return nil,201,""
}

func Dequeue(strs []string)(error,int,string){
  if len(myQueue[strs[0]]) == 0{
    return errors.New("queue is empty"),400,""
  }
  a := myQueue[strs[0]][0]
  myQueue[strs[0]] = myQueue[strs[0]][1:]

  fmt.Println(myQueue[strs[0]])

  return nil,200,strconv.Itoa(a)
}

func ProcessCommand(strs []string) (error,int,string){
  command := strs[0]
  if command == "SET"{
   return HandleSetCommand(strs)
  }

  if command == "GET"{
    key := strs[1]
    if myMap[key].exists && time.Now().Before(myMap[key].expires){
      return nil,200,myMap[key].value
    }else{
      return errors.New("key not found"),404,""
    }
  }

  if command == "QPUSH"{
    return Enqueue(strs[1:])
  }

  if command == "QPOP"{
    return Dequeue(strs[1:])
  }

  return errors.New("invalid command"),400,""
}

func checkExistence(a,condition string) error{
  if condition == "NX"{
    if myMap[a].exists {
      return errors.New("key already exist")
    }
  }else if condition == "XX" {
    if !myMap[a].exists {
      return errors.New("key not found")
    }
  }else{
    return errors.New("invalid argument")
  }

  return nil
}
