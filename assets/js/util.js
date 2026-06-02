export function pad(number, digits) {
    let s = "000000000" + number;
    return s.substring(s.length-digits);
}