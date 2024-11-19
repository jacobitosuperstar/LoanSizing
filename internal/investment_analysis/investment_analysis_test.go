package investment_analysis
import (
  "testing";
  "fmt";
  ls "jacobitosuperstar/LoanSizing/internal/loan_sizer";
)


func TestROI(t *testing.T) {
    var testCases = []struct {
        input  ReturnOfInvestment
        want float64
    }{
      {
        ReturnOfInvestment{
          taxMetrics: TaxAssumptions{
            LanBuildingValue: 0.3,
            FixDepreciationTimeLine: 27,
            IncomeTaxRate: 0.25,
            CapitalGainsTaxRate: 0.15,
            DepreciationRecaptureTaxRate: 0.25,
          },
          dealMetrics: DealInformation{
            PurchasePrice: 6500000,
            ClosingAndRenovations: 225000,
            GoingInCapRate: 0.0596,
            InitRevenue: 687500,
            InitOperatingExpenses: 300000,
            InitCapitalReserves: 7500,
            ProjRevenueGrowth: 0.0350,
            ProjOperatingExpensesGrowth: 0.0250,
            ProjCapitalReservesGrowth: 0.0250,
          },
          loanMetrics: ls.LoanSizer{
            MaxLTV: 0.70,
            MinDSCR: 1.25,
            Amortization: 30,
            Term: 10,
            IOPeriod: 2,
            Rate: 0.0450,
            PropertyValue: 0.0,
            NOI: 0.0,
            RequestedLoanAmount: 1000000000,
            LoanOriginationFees: 0.01,
          },
          saleMetrics: SaleTerms{
            ExitCapRate: 0.0650,
            CostOfSale: 0.0250,
            SaleYear: 10,
          },
        },
        0.0,
      },
    }

    for _, test := range testCases {
      roi, err := InitTargetReturnOfInvestment(test.input)
      if err != nil {
          t.Errorf("got: %v, error: %v", roi, err)
      }
      fmt.Printf("roi: %+v\n", roi)
      // t.Errorf("got: %v", got)
    }
}
