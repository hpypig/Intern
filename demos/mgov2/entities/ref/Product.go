package ref

import (
    "fmt"
    "gopkg.in/mgo.v2/bson"
    "strings"
    "time"
)

type Product struct {
    Id         bson.ObjectId `bson:"_id"`
    Name       string        `bson:"name"`
    Price      float64       `bson:"price"`
    Quantity   int64         `bson:"quantity"`
    Status     bool          `bson:"status"`
    Date       time.Time     `bson:"date"`
    CategoryId bson.ObjectId `bson:"categoryId"`
    Brand      Brand         `bson:"brand"`
    Colors     []string      `bson:"colors"`
}

func (this Product) ToString() string {
    result := fmt.Sprintf("id: %s", this.Id.Hex())
    result = result + fmt.Sprintf("\nname: %s", this.Name)
    result = result + fmt.Sprintf("\nprice: %0.1f", this.Price)
    result = result + fmt.Sprintf("\nquantity: %d", this.Quantity)
    result = result + fmt.Sprintf("\nstatus: %t", this.Status)
    result = result + fmt.Sprintf("\ndate: %s", this.Date.Format("2006-01-02"))
    //result = result + fmt.Sprintf("\ncategory id: %s", this.Category.Hex())
    result = result + this.Brand.ToString()
    result = result + fmt.Sprintf("\ncolors: %s", strings.Join(this.Colors, ", "))
    return result
}
