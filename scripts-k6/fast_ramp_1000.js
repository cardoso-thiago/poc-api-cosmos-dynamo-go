import http from 'k6/http';
import { SharedArray } from 'k6/data';

export const options = {
  scenarios: {
    steady_load: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s',
      duration: '3m',
      preAllocatedVUs: 200,
      maxVUs: 2000,
    },
  },
};

const uuids = new SharedArray('uuids', function() {
    return open('../uuids.txt').split('\n').filter(u => u.trim().length > 0);
});

export default function () {
    const uuid = uuids[Math.floor(Math.random() * uuids.length)];
    const res = http.get(`http://localhost:8888/items/${uuid}`);
}