import { test } from 'node:test'
import assert from 'node:assert/strict'
import {pad} from "../../assets/js/util.js";

test('pad test', ()=> {
    let rez = pad( 1, 2)

    assert.equal(rez, '01')
})