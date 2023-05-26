package ref
import (
    "fmt"

    "gopkg.in/mgo.v2/bson"
)

type CategoryGroup struct {
    CategoryId   bson.ObjectId
    CountProduct int
    SumQuantity  int64
    MinPrice     float64
    MaxPrice     float64
    AvgPrice     float64
}

func (this CategoryGroup) ToString() string {
    result := fmt.Sprintf("Category Id: %s\n", this.CategoryId.Hex())
    result = result + fmt.Sprintf("Count Product: %d\n", this.CountProduct)
    result = result + fmt.Sprintf("Sum Quantities: %d\n", this.SumQuantity)
    result = result + fmt.Sprintf("Min Price: %0.1f\n", this.MinPrice)
    result = result + fmt.Sprintf("Max Price: %0.1f", this.MaxPrice)
    return result
}
