// The main idea here is do as the package name says investment analysis in
// our real state buys

// TODO: how is the projected return of our investment?
// [ ] TaxAssumptions
// [ ] DealMetrics
// [ ] LoanTerms
// [ ] SaleMetrics
// [ ] ReturnOfInvestment

package investment_analysis

import (
    "fmt";
    "math";
    ff "jacobitosuperstar/LoanSizing/internal/financial_formulas";
    ls "jacobitosuperstar/LoanSizing/internal/loan_sizer";
)

// TaxAssumptions is a struct that has all the taxes information regarding the
// deal.
type TaxAssumptions struct {
    LanBuildingValue                float64     `json:"lan_building_value"`
    FixDepreciationTimeLine         int         `json:"fixed_depreciation_timeline"`
    IncomeTaxRate                   float64     `json:"income_tax_rate"`
    CapitalGainsTaxRate             float64     `json:"capital_gains_tax_rate"`
    DepreciationRecaptureTaxRate    float64     `json:"depreciation_recapture_tax_rate"`
}

// DealInformation is a struct that has all the information regarding the
// buying of the commercial property.
type DealInformation struct {
    PurchasePrice               int         `json:"purchase_price"`
    ClosingAndRenovations       int         `json:"closing_and_renovations"`
    GoingInCapRate              float64     `json:"going_in_caprate"`
    InitRevenue                 float64     `json:"initial_revenue"`
    InitOperatingExpenses       float64     `json:"initial_operating_expenses"`
    InitCapitalReserves         float64     `json:"initial_capital_reserves"`
    ProjRevenueGrowth           float64     `json:"projected_revenue_growth"`
    ProjOperatingExpensesGrowth float64     `json:"projected_operating_expenses_growth"`
    ProjCapitalReservesGrowth   float64     `json:"projected_capital_reserves_growth"`
}

// SaleTerms is a struc that has all the sale information regarding the sale of
// the sale of the property.
type SaleTerms struct {
    ExitCapRate         float64     `json:"exit_cap_rate"`
    CostOfSale          float64     `json:"cost_of_sale"`
    SaleYear            int         `json:"sale_year"`
}

// ProjectedSalePrice returns the projected sale price of real state.
func (st SaleTerms) ProjectedSalePrice (noi float64) float64 {
    projected_sale_price := ff.Round2(noi / st.ExitCapRate)
    projected_sale_price = projected_sale_price - (projected_sale_price * st.CostOfSale)
    return ff.Round2(projected_sale_price)
}

// ROI of the totallity of the deal.
type ReturnOfInvestment struct {
    taxMetrics              TaxAssumptions
    dealMetrics             DealInformation
    loanMetrics             ls.LoanSizer
    saleMetrics             SaleTerms
    // Calculated fields
    AdquisitionCost         float64                     `json:"adquisition_cost"`
    NetCashFlowProjection   []map[string]interface{}    `json:"net_cash_flow_projection"`
    IRR                     float64                     `json:"internal_rate_of_return"`
    EquityMultiple          float64                     `json:"equity_multiple"`
    AverageCashOnCashReturn float64                     `json:"average_cash_on_cash_return"`
}

// SetAdquisitionCost sets the AdquisitionCost of Deal
func (roi *ReturnOfInvestment) SetAdquisitionCost ()  {
    adquisitionCost := - float64(roi.dealMetrics.PurchasePrice) -
    float64(roi.dealMetrics.ClosingAndRenovations) -
    (roi.loanMetrics.LoanOriginationFees * roi.loanMetrics.MaximumLoanAmount) +
    roi.loanMetrics.MaximumLoanAmount
    roi.AdquisitionCost = adquisitionCost
}

// CashOnCashReturn returns the made money in reference to the money invested
// to adquire the property.
func (roi ReturnOfInvestment) CashOnCashReturn (net_cash_flow float64)  float64 {
    return ff.Round4(math.Abs(net_cash_flow/roi.AdquisitionCost))
}

// SetNetCashFlowProjection sets the NetCashFlowProjection of the Deal
func (roi *ReturnOfInvestment) SetNetCashFlowProjection () error {
    var net_cash_flow_projection []map[string]interface{}
    net_cash_flow_projection = append(
        net_cash_flow_projection,
        map[string]interface{} {
            "net_cash_flow": roi.AdquisitionCost,
        },
    )

    revenue := roi.dealMetrics.InitRevenue
    expense := roi.dealMetrics.InitOperatingExpenses
    reserve := roi.dealMetrics.InitCapitalReserves

    // getting the building value
    purchase_price := roi.dealMetrics.PurchasePrice
    building_value := ff.Round2(float64(purchase_price) * (1.0 - roi.taxMetrics.LanBuildingValue))

    // depreciation of the building
    building_depreciation := ff.Round2(- building_value/float64(roi.taxMetrics.FixDepreciationTimeLine))

    // payment distribution of the loan
    ppmt, ipmt, err := roi.loanMetrics.PaymentDistribution()

    if err != nil {
        return fmt.Errorf("PaymentDistribution internal error: %v", err)
    }

    for i := 1; i <= roi.loanMetrics.Term; i++ {
        // this year NOI
        current_noi := ff.Round2(revenue - expense)
        // this year interest and principal payments
        current_ppmt := ppmt[i-1]
        current_ipmt := ipmt[i-1]
        current_pmt := roi.loanMetrics.LoanPayment
        // cashflow after debt service
        cfads := ff.Round2(current_noi + reserve + current_pmt)
        // depreciation expense
        depreciation_expense := 0.0
        if i < roi.taxMetrics.FixDepreciationTimeLine {
            depreciation_expense = building_depreciation
        }
        // income tax
        income_tax := ff.Round2(- (current_noi + current_ipmt + depreciation_expense) * roi.taxMetrics.IncomeTaxRate)
        implied_income_tax := ff.Round4(math.Abs(income_tax/cfads))
        // net cashflow
        ncf := ff.Round2(cfads + income_tax)
        // cash on cash return
        cocr := roi.CashOnCashReturn(ncf)

        net_cash_flow_projection = append(
            net_cash_flow_projection,
            map[string]interface{} {
                "year": i,
                "revenue": revenue,
                "expense": expense,
                "noi": current_noi,
                "reserve": reserve,
                "principal_payment": current_ppmt,
                "interest_payment": current_ipmt,
                "cashflow_after_debt_service": cfads,
                "depreciation_expense": depreciation_expense,
                "income_tax": income_tax,
                "implied_income_tax": implied_income_tax,
                "net_cash_flow": ncf,
                "cash_on_cash_return": cocr,
            },
        )
        revenue += ff.Round2(revenue * roi.dealMetrics.ProjRevenueGrowth)
        expense += ff.Round2(expense * roi.dealMetrics.ProjOperatingExpensesGrowth)
        reserve += ff.Round2(reserve * roi.dealMetrics.ProjCapitalReservesGrowth)
    }
    after_term_noi := ff.Round2(revenue + expense)
    // Adding the cashflow after the sell of the property
    // sale with the projected NOI
    projected_sale_price := roi.saleMetrics.ProjectedSalePrice(after_term_noi)
    // capital gains tax
    cg := projected_sale_price -
        float64(roi.dealMetrics.PurchasePrice) -
        float64(roi.dealMetrics.ClosingAndRenovations)
    cg = ff.Round2(cg)
    cgt := cg * roi.taxMetrics.CapitalGainsTaxRate
    // Depreciation Recapture tax
    drt := building_depreciation *
        float64(roi.saleMetrics.SaleYear) *
        roi.taxMetrics.DepreciationRecaptureTaxRate
    drt = ff.Round2(drt)
    // Sale calculations
    sale := net_cash_flow_projection[roi.loanMetrics.Term - 1]
    sale_net_cash_flow := sale["net_cash_flow"].(float64)
    sale_net_cash_flow = sale_net_cash_flow +
        projected_sale_price +
        drt +
        cgt +
        roi.loanMetrics.BalloonPayment
    sale_net_cash_flow = ff.Round2(sale_net_cash_flow)
    sale["net_cash_flow"] = sale_net_cash_flow
    sale["sale_price"] = projected_sale_price
    sale["depreciation_recapture_tax"] = drt
    sale["capital_gains_tax"] = cgt
    // Setting the value
    roi.NetCashFlowProjection = net_cash_flow_projection
    return nil
}

func (roi *ReturnOfInvestment) SetIRR () error {
    return nil
}

func (roi *ReturnOfInvestment) SetEquityMultiple () error {
    return nil
}

func (roi *ReturnOfInvestment) SetAverageCashOnCashReturn () error {
    return nil
}

// InitReturnOfInvestment sets the calculated terms in the ReturnOfInvestment
// struct.
func InitReturnOfInvestment(roi ReturnOfInvestment) ReturnOfInvestment {
    return roi
}

// Target ROI of the investment
type TargetReturnOfInvestment struct {
    IRR                     float64     `json:"internal_rate_of_return"`
    EquityMultiple          float64     `json:"equity_multiple"`
    AverageCashOnCashReturn float64     `json:"average_cash_on_cash_return"`
}

//
func InitTargetReturnOfInvestment (roi ReturnOfInvestment)  (ReturnOfInvestment, error) {
    roi.SetAdquisitionCost()
    err := roi.SetNetCashFlowProjection()
    if err != nil {
        return roi, err
    }
    return roi, nil
}
