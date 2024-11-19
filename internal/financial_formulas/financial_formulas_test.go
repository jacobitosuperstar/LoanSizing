// Testing of the Financial Formulas
// TODO: Tests needed
// [X] Round2
// [X] IOPayment (monthly)
// [X] YearlyIOPayment
// [X] Payment (monthly)
// [X] YearlyPayment
// [X] PrincipalPayment (monthly)
// [X] InterestPayment (monthly)
// [X] PresentValue

package financial_formulas
import "testing"

func TestRound2(t *testing.T) {
    var testCases = []struct {
        input float64
        want float64
    }{
        {0.000, 0.00},
        {0.080, 0.08},
        {1.994, 1.99},
        {1.995, 2.00},
        {1.999, 2.00},
    }

    for _, test := range testCases {
        if got := Round2(test.input); got != test.want {
            t.Errorf("got: %g, wanted: %g", got, test.want)
        }
    }
}

func TestIOPayment(t *testing.T){
    var testCases = []struct {
        rate float64
        pv float64
        want float64
    }{
        {0.00, 0.00, 0.00},
        {100, 0.00375, - 0.38},
        {200, 0.00375, - 0.75},
    }

    for _, test := range testCases {
        if got  := IOPayment(test.rate, test.pv); got != test.want {
            t.Errorf("got: %g, wanted: %g", got, test.want)
        }
    }
}

// func TestYearlyIOPayment(t *testing.T){
//     var testCases = []struct {
//         rate float64
//         pv float64
//         want float64
//     }{
//         {0.00, 0.00, 0.00},
//         {100, 0.00375, 4.50},
//         {200, 0.00375, 9.00},
//     }
//
//     for _, test := range testCases {
//         if got, _ := YearlyIOPayment(test.rate, test.pv); got != test.want {
//             t.Errorf("got: %g, wanted: %g", got, test.want)
//         }
//     }
// }

func TestPayment(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pv float64
        fv float64
        paymentType int
        want float64
    }{
        {
            name: "Everything Zero",
            rate: 0,
            numPeriods: 0,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Valid numPeriods",
            rate: 0,
            numPeriods: 12,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Invalid paymentType",
            rate: 0.00375,
            numPeriods: 12,
            pv: 100,
            fv: 0,
            paymentType: 5,
            want: 0,
        },
        {
            name: "Totally valid case",
            rate: 0.00375,
            numPeriods: 12,
            pv: 100,
            fv: 0,
            paymentType: 0,
            want: -8.54,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, _ := Payment(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)
            if got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}

// func TestYearlyPayment(t *testing.T){
//     var testCases = []struct {
//         name string
//         rate float64
//         numPeriods int
//         pv float64
//         fv float64
//         paymentType int
//         want float64
//     }{
//         {
//             name: "Everything Zero",
//             rate: 0,
//             numPeriods: 0,
//             pv: 0,
//             fv: 0,
//             paymentType: 0,
//             want: 0,
//         },
//         {
//             name: "Everything Zero, Valid numPeriods",
//             rate: 0,
//             numPeriods: 12,
//             pv: 0,
//             fv: 0,
//             paymentType: 0,
//             want: 0,
//         },
//         {
//             name: "Invalid paymentType",
//             rate: 0.00375,
//             numPeriods: 12,
//             pv: 100,
//             fv: 0,
//             paymentType: 5,
//             want: 0,
//         },
//         {
//             name: "Totally valid case",
//             rate: 0.00375,
//             numPeriods: 12,
//             pv: 100,
//             fv: 0,
//             paymentType: 0,
//             want: -102.48,
//         },
//     }
//
//     for _, test := range testCases {
//         t.Run(test.name, func(t *testing.T) {
//             got, _ := YearlyPayment(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)
//             if got != test.want {
//                 t.Errorf("got: %g, wanted: %g", got, test.want)
//             }
//         })
//     }
// }

func TestPrincipalPayments(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pv float64
        fv float64
        paymentType int
        want []float64
    }{
        {
            name: "Everything Zero",
            rate: 0,
            numPeriods: 0,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: []float64{},
        },
        {
            name: "Everything Zero, Valid numPeriods",
            rate: 0,
            numPeriods: 1,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: []float64{0},
        },
        {
            name: "Invalid paymentType",
            rate: 0,
            numPeriods: 1,
            pv: 0,
            fv: 0,
            paymentType: 5,
            want: []float64{},
        },
        {
            name: "Totally valid case",
            rate: 0.00375,
            numPeriods: 2,
            pv: 100,
            fv: 0,
            paymentType: 0,
            want: []float64{-49.9, -50.09},
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, _ := PrincipalPayments(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)

            if len(got) != len(test.want) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            } else {
                for i := range got {
                    if got[i] != test.want[i] {
                        t.Errorf("got: %g, wanted: %g", got[i], test.want[i])
                    }
                }
            }
        })
    }
}

func TestInterestPayments(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pv float64
        fv float64
        paymentType int
        want []float64
    }{
        {
            name: "Everything Zero",
            rate: 0,
            numPeriods: 0,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: []float64{},
        },
        {
            name: "Everything Zero, Valid numPeriods",
            rate: 0,
            numPeriods: 1,
            pv: 0,
            fv: 0,
            paymentType: 0,
            want: []float64{0},
        },
        {
            name: "Invalid paymentType",
            rate: 0,
            numPeriods: 1,
            pv: 0,
            fv: 0,
            paymentType: 5,
            want: []float64{},
        },
        {
            name: "Totally valid case",
            rate: 0.00375,
            numPeriods: 2,
            pv: 100,
            fv: 0,
            paymentType: 0,
            want: []float64{-0.38, -0.19},
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, _ := InterestPayments(test.rate, test.numPeriods, test.pv, test.fv, test.paymentType)

            if len(got) != len(test.want) {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            } else {
                for i := range got {
                    if got[i] != test.want[i] {
                        t.Errorf("got: %g, wanted: %g", got[i], test.want[i])
                    }
                }
            }
        })
    }
}

func TestPresentValue(t *testing.T){
    var testCases = []struct {
        name string
        rate float64
        numPeriods int
        pmt float64
        fv float64
        paymentType int
        want float64
    }{
        {
            name: "Everything Zero",
            rate: 0,
            numPeriods: 0,
            pmt: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Everything Zero, Valid numPeriods",
            rate: 0,
            numPeriods: 1,
            pmt: 0,
            fv: 0,
            paymentType: 0,
            want: 0,
        },
        {
            name: "Invalid paymentType",
            rate: 0,
            numPeriods: 1,
            pmt: 0,
            fv: 0,
            paymentType: 5,
            want: 0,
        },
        {
            name: "Totally valid case",
            rate: 0.08,
            numPeriods: 20,
            pmt: 500,
            fv: 0,
            paymentType: 0,
            want: -4909.07,
        },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            got, _ := PresentValue(test.rate, test.numPeriods, test.pmt, test.fv, test.paymentType)
            if got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}
