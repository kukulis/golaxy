import {pad} from "./util.js";

export function formatDateTime(value) {
    if (!value) return ''
    const date = new Date(value);

    // console.log( 'date object is ', date )

    if ( !(date instanceof  Date) || isNaN(date) ) {
        console.log('Invalid date detected')
        return '';
    }

    let y = date.getFullYear()
    let M = pad((date.getMonth() + 1), 2)
    let d = pad(date.getDate(), 2)

    let h = pad(date.getHours(), 2)
    let m = pad(date.getMinutes(), 2)
    let s = pad(date.getSeconds(), 2)

    // let z = date.getTimezoneOffset()

    return `${y}-${M}-${d} ${h}:${m}:${s}`
}
