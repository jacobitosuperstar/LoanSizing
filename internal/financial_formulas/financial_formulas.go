// Opinionated Financial Formulas. Every money value rounded down to 2 decimal
// places.
// TODO: Formulas needed
// [X] IOPayment (monthly)
// [X] YearlyIOPayment
// [X] Payment (monthly)
// [X] YearlyPayment
// [X] PrincipalPayment (monthly)
// [X] InterestPayment (monthly)
// [X] PresentValue

package financial_formulas

import (
	"errors"
	"math"
)

const (
    PayEnd = iota
    PayBegin
)


// round2 returns a float number rounded to 2 decimals
func round2(num float64) float64 {
    return (math.Round(num*100)/100)
}

// IOPayment returns the interest only payment for a cash flow with a constant
// interest rate.
func IOPayment(
    rate float64,
    pv float64,
) (
    pmt float64,
    err error,
) {
    pmt = (pv*rate)
    return round2(pmt), nil
}

// YearlyIOPayment returns the yearly interest only payment for a cash flow
// with a constant interest rate.
func YearlyIOPayment(
    rate float64,
    pv float64,
) (
    pmt float64,
    err error,
) {
    pmt = 12*(pv*rate)
    return round2(pmt), nil
}

// Payment returns the constant payment for a cash flow with a constant
// interest rate.
func Payment(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    pmt float64,
    err error,
) {
    if numPeriods <= 0 {
		return 0.0, errors.New("The periods must be greater than 0")
	}
	if paymentType != PayEnd && paymentType != PayBegin {
		return 0.0, errors.New("payment type must be pay-end or pay-begin")
	}
	if rate != 0 {
		pmt = (-fv - pv*math.Pow(1+rate, float64(numPeriods))) / (1 + rate*float64(paymentType)) / ((math.Pow(1+rate, float64(numPeriods)) - 1) / rate)
	} else {
		pmt = (-pv - fv) / float64(numPeriods)
	}
	return round2(pmt), nil
}

// YearlyPayment returns the yearly loan payment for a cash flow with a
// constant interest rate.
func YearlyPayment(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    pmt float64,
    err error,
) {
    pmt, err = Payment(rate, numPeriods, pv, fv, paymentType)

    if err != nil {
        return 0, err
    }
    return round2(12*pmt), nil
}

// InterestAndPrincipalPayment returns an array of interest payments, an
// principal payments and possibly an error.
func InterestAndPrincipalPayments(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    ipmt []float64,
    ppmt []float64,
    err error,
) {
    pmt, err := Payment(rate, numPeriods, pv, fv, paymentType)

    if err != nil {
        return ipmt, ppmt, err
    }

    capital := pv

    for i := 1; i <= numPeriods; i++ {
        if paymentType == PayBegin && i == 1 {
            ipmt = append(ipmt, 0.00)
        } else {
            interest_payment := round2(- capital * rate)
            ipmt = append(ipmt, interest_payment)

            principal_payment := round2(pmt - interest_payment)
            ppmt = append(ppmt, principal_payment)

            capital += principal_payment
        }
    }

    if capital != fv {
        return ipmt, ppmt, errors.New("The final values of the loan do not match")
    }
    return ipmt, ppmt, nil
}

// PrincipalPayments return an array and an error of all the principal payments
// during the number of periods
func PrincipalPayments(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    []float64,
    error,
) {
    _, ppmt, err := InterestAndPrincipalPayments(rate, numPeriods, pv, fv, paymentType)

    if err != nil {
        return ppmt, err
    }
    return ppmt, nil
}

// InterestPayments returns an array and an error of all the interest payments
// during the number of periods
func InterestPayments(
    rate float64,
    numPeriods int,
    pv float64,
    fv float64,
    paymentType int,
) (
    []float64,
    error,
) {
    ipmt, _, err := InterestAndPrincipalPayments(rate, numPeriods, pv, fv, paymentType)

    if err != nil {
        return ipmt, err
    }
    return ipmt, nil
}

// PresentValue return the present value of a cashflow with constant interest
// rate and payments
func PresentValue(
    rate float64,
    numPeriods int,
    pmt float64,
    fv float64,
    paymentType int,
) (
    pv float64,
    err error,
) {
    if numPeriods <= 0 {
        return 0, errors.New("Number of periods must be greater than 0")
    }
    if paymentType != PayEnd && paymentType != PayBegin {
        return 0, errors.New("Payment must be pay-end or pay-begin")
    }
    if rate != 0 {
        pv = (-pmt*(1+rate*float64(paymentType))*((math.Pow(1+rate, float64(numPeriods))-1)/rate) - fv) / math.Pow(1+rate, float64(numPeriods))
    } else {
        pv = -fv - pmt*float64(numPeriods)
    }
    return round2(pv), nil
}
