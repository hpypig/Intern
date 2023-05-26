package entities
type BaseParam struct {
    Page int64 `form:"page" binding:"required"`
    PageSize int64 `form:"pageSize" binding:"required"`
}
type ParamStockNews struct { // 这个个也可以改成 request
    MarketCode string `form:"ids" binding:"required"`
    TxtType int `form:"txtType" binding:"required"`
    BaseParam
}
type ParamColumnRequest struct {
    Id string `form:"id" binding:"required"`
    BaseParam
}


func NewParamStockNews() *ParamStockNews {
    return &ParamStockNews {
        BaseParam: BaseParam{
            Page:1,
            PageSize:10,
        },
    }
}

func NewParamColumnRequest() *ParamColumnRequest {
    return &ParamColumnRequest {
        BaseParam: BaseParam{
            Page:1,
            PageSize:10,
        },
    }
}
