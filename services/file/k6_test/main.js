import http, { file } from 'k6/http';
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.2/index.js';

export const options = {
    thresholds: {
        http_req_duration: [{
            threshold: 'avg<2000', // Average latency below 2s
            abortOnFail: true,
            delayAbortEval: '10s'
        }],
    },
    scenarios: {
        average_load: {
            executor: 'ramping-vus',
            stages: [
                { duration: '10s', target: 20 },
                { duration: '10s', target: 120 },
                { duration: '10s', target: 300 },
            ]
        }
    },
};

const smallFile = open('./figure/image-50KB.jpg', 'b');
const medFile = open('./figure/image-90KB.jpg', 'b');
const bigFile = open('./figure/image-100KB.jpg', 'b');
const hugeFile = open('./figure/image-200KB.jpg', 'b');
const invalidFile = open('./figure/sql-5KB.sql', 'b');

export default function() {
    const config = {
        baseUrl: __ENV.BASE_URL ? __ENV.BASE_URL : "http://localhost:8080"
    }
    var fileToTest = {
        small: smallFile,
        smallName: "small.jpg",
        medium: medFile,
        mediumName: "med.jpg",
        big: bigFile,
        bigName: "big.jpg",
        huge: hugeFile,
        hugeName:  "huge.jpg",
        invalid: invalidFile,
        invalidName: "invalid.sql",
    }
    const route = config.baseUrl + "/v1/file";
    describe('file service: positive payload', () => {
        const positivePayloads = [
            { file: file(fileToTest.small, fileToTest.smallName) },
            { file: file(fileToTest.medium, fileToTest.mediumName) },
        ]
        positivePayloads.forEach((payload) => {
            const res = http.post(route, payload)
            expect(res.status, 'response status').to.equal(200);
            expect(res).to.have.validJsonBody();
        })
    })
    describe('file service: negative payload', () => {
        describe('invalid file size', () => {
            const negativePayloads = [
                { file: file(fileToTest.big, fileToTest.bigName) },
                { file: file(fileToTest.huge, fileToTest.hugeName) },
            ]
            negativePayloads.forEach((payload) => {
                const res = http.post(route, payload)
                expect(res.status, 'response status').to.equal(400);
            })
        })
        describe('invalid file type', () => {
            const negativePayloads = [
                { file: file(fileToTest.invalid, fileToTest.invalid) },
            ]
            negativePayloads.forEach((payload) => {
                const res = http.post(route, payload)
                expect(res.status, 'response status').to.equal(400);
            })
        })
    })
}