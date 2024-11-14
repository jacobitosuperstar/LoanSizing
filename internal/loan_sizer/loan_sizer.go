// We will create a LoanSizer struct and we will be adding things to the struct
// as is needed to show the different things needed or wanted from the loan
// sizing part.
// TODO: Things that should be returned to the user
// [X] maximum loan amount
// [X] requested loan amount
// [X] yearly loan payment
// [X] yearly loan io period payment
// [X] balloon payment

package loan_sizer

import (
    "fmt"
    // "log"
    "math"
    "sort"
    ff "jacobitosuperstar/LoanSizing/internal/financial_formulas";
)

type LoanSizer struct {
    MaxLTV float64                  `json:"max_ltv"`
    MinDSCR float64                 `json:"min_dscr"`
    Amortization int                `json:"amortization"`
    Term int                        `json:"term"`
    IOPeriod int                    `json:"io_period"`
    Rate float64                    `json:"rate"`
    PropertyValue int               `json:"property_value"`
    NOI float64                     `json:"noi"`
    RequestedLoanAmount int         `json:"requested_loan_amount"`
    LoanOriginationFees float64     `json:"loan_origination_fees"`
    // Calculated fields
    MaximumLoanAmount float64       `json:"maximum_loan_amount"`
    LoanPayment float64             `json:"yearly_loan_payment"`
    IOLoanPayment float64           `json:"yearly_io_loan_payment"`
    BalloonPayment float64          `json:"balloon_payment"`
}

// Calculation methods

// max_ltv_loan_amount returns the maximum loan amount given the maximum loan
// to value ratio
func (ls LoanSizer) max_ltv_loan_amount () float64 {
    ltv_mla  := math.Floor(ls.MaxLTV * float64(ls.PropertyValue))
    return ltv_mla
}

// max_mindscr_loan_amount returns the maximum loan amount given the minimum
// dscr
func (ls LoanSizer) max_mindscr_loan_amount () (float64, error) {
    // monthly_rate := ls.Rate / 12
    // amoritzation_months := ls.Amortization * 12
    payment := - ls.NOI / ls.MinDSCR
    dscr_mla, err := ff.PresentValue(ls.Rate, ls.Amortization, payment, 0, 0)

    if err != nil {
        return 0.0, fmt.Errorf("max_mindscr_loan_amount internal error: %v", err)
    }
    return math.Floor(dscr_mla), err
}

// Setter methods

// SetMaximumLoanAmount sets the maximum loan amount of a LoanSizer struct
func (ls *LoanSizer) SetMaximumLoanAmount () error {
    max_mindscr_loan_amount, err := ls.max_mindscr_loan_amount()

    if err != nil {
        ls.MaximumLoanAmount = 0.0
        return fmt.Errorf("ls.max_mindscr_loan_amount internal error: %v", err)
    }

    loan_values := [3]float64{
        ls.max_ltv_loan_amount(),
        max_mindscr_loan_amount,
        float64(ls.RequestedLoanAmount),
    }
    sort.Float64s(loan_values[:])
    ls.MaximumLoanAmount = loan_values[0]
    return nil
}

// SetIOLoanPayment sets the loan payments during the IO periods
func (ls *LoanSizer) SetIOLoanPayment () {
    ls.IOLoanPayment = ff.IOPayment(ls.Rate, ls.MaximumLoanAmount)
}

// SetLoanPayment sets the yearly loan payments for the maximum amount
func (ls *LoanSizer) SetLoanPayment () error {
    loan_payment, err := ff.Payment(ls.Rate, ls.Amortization, ls.MaximumLoanAmount, 0, 0)
    if err != nil {
        return fmt.Errorf("Payment internal error: %v", err)
    }
    ls.LoanPayment = loan_payment
    return nil
}

// SetBallonPayment sets the balloon payment at the end of the term
func (ls *LoanSizer) SetBalloonPayment () error {
    principal_payments, err := ff.PrincipalPayments(ls.Rate, ls.Amortization, ls.MaximumLoanAmount, 0, 0)

    if err != nil {
        ls.BalloonPayment = 0.0
        return fmt.Errorf("PrincipalPayments internal error: %v", err)
    }

    if ls.IOPeriod > 0 {
        io_period_ppmt := make([]float64, ls.IOPeriod)
        principal_payments = append(io_period_ppmt, principal_payments...)
    }
    // There is no need to create a new slice that will contain only the term
    // as the iteration will iterate till the term value. Genious move!!.

    capital := ls.MaximumLoanAmount
    for i:=0; i < ls.Term; i++ {
        capital += principal_payments[i]
    }
    capital = ff.Round2(capital)
    ls.BalloonPayment = capital
    return nil
}
