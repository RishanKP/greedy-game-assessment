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

func checkExistence(a, condition string) error {
	if condition == "NX" {
		if myMap[a].exists {
			return errors.New("key already exist")
		}
	} else if condition == "XX" {
		if !myMap[a].exists {
			return errors.New("key not found")
		}
	}
	return nil
}

func HandleSetCommand(strs []string) (error, int, string) {
	key := strs[0]
	value, err := strconv.Atoi(strs[1])
	if err != nil {
		return errors.New("invalid value given. expects integer"), 400, ""
	}
	expirySet := false
	var expires time.Time

	for i := 2; i < len(strs); i++ {
		if strs[i] == "EX" {
			duration, e := strconv.Atoi(strs[i+1])
			if e != nil {
				return errors.New("invalid expiry duration"), 400, ""
			}
			expires = time.Now().Add(time.Second * time.Duration(duration))
			expirySet = true

			i = i + 1
		}

		if strs[i] == "NX" || strs[i] == "XX" {
			err = checkExistence(key, strs[i])
			if err != nil {
				return err, 400, ""
			}
		}
	}

	if expirySet == false {
		expires = time.Now().AddDate(20, 0, 0)
	}

	myMap[key] = mydb{value: value, expires: expires, exists: true}

	fmt.Println(myMap[key].expires)
	return nil, 201, ""
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
    return HandleSetCommand(strs[1:])
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
