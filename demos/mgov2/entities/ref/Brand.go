package ref
import (
    "fmt"
    "gopkg.in/mgo.v2/bson"
)

//type Brand struct {
//    Id   bson.ObjectId `bson:"_id"`
//    Name string        `bson:"name"`
//}

type Brand struct {
    Id   bson.ObjectId
    Name string
}

func (this Brand) ToString() string {
    result := fmt.Sprintf("\nbrand id: %s", this.Id.Hex())
    result = result + fmt.Sprintf("\nbrand name: %s", this.Name)
    return result
}
