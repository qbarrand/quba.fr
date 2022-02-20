// import {Constraint, generateMediaConstraints, LANDSCAPE, Orientation, PORTRAIT} from "./constraint";
//
// const cases = [
//     [PORTRAIT, 0, 1000, '(orientation: portrait) and (max-height: 1000px)'],
//     [PORTRAIT, 1000, 2000, '(orientation: portrait) and (min-height: 1000px) and (max-height: 2000px)'],
//     [LANDSCAPE, 1000, Infinity, '(orientation: landscape) and (min-width: 1000px)'],
//     [LANDSCAPE, 1000, 2000, '(orientation: landscape) and (min-width: 1000px) and (max-width: 2000px)'],
// ]
//
// test.each(cases)('%s toCSSMediaQuery should return %s', (o: Orientation, l:number, h: number, expected: string) => {
//     expect(
//         new Constraint(o, l, h).toMediaQuery()
//     ).toBe(
//         expected
//     )
// })
//
// test('2 widths, 2 heights', () => {
//     expect(
//         generateMediaConstraints([10, 20], [30, 40])
//     ).toEqual(
//         {
//             '(orientation: landscape) and (max-width: 10px)': new Constraint(LANDSCAPE, 0, 10),
//             '(orientation: landscape) and (min-width: 11px) and (max-width: 20px)': new Constraint(LANDSCAPE, 11, 20),
//             '(orientation: landscape) and (min-width: 21px)': new Constraint(LANDSCAPE, 21, Infinity),
//             '(orientation: portrait) and (max-height: 30px)': new Constraint(PORTRAIT, 0, 30),
//             '(orientation: portrait) and (min-height: 31px) and (max-height: 40px)': new Constraint(PORTRAIT, 31, 40),
//             '(orientation: portrait) and (min-height: 41px)': new Constraint(PORTRAIT, 41, Infinity),
//         }
//     )
// })