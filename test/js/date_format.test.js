import { test } from 'node:test'
import assert from 'node:assert/strict'
import { formatDateTime } from '../../assets/js/date_format.js'

test('formatDateTime formats ISO timestamp', () => {
    const result = formatDateTime('2026-05-07T05:04:08.840119403+03:00')
    assert.equal( result, '2026-05-07 05:04:08' )
})

test('formatDateTime invalid date', () => {
    const result = formatDateTime('2026-05-071T05:04:08.840119403+03:00')
    assert.equal( result, '' )
})

test('formatDateTime returns empty string for null', () => {
    assert.equal(formatDateTime(null), '')
})

test('formatDateTime returns empty string for empty string', () => {
    assert.equal(formatDateTime(''), '')
})

