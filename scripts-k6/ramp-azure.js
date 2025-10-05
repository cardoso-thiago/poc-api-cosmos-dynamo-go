import http from 'k6/http';
import { SharedArray } from 'k6/data';

export const options = {
  scenarios: {
    ramping_load: {
      executor: 'ramping-arrival-rate',
      startRate: 500,
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 2000,
      stages: [
        { target: 1000, duration: '1m' },
        { target: 1500, duration: '3m' },
      ],
    },
  },
};

const uuids = new SharedArray('uuids', function() {
    return open('../docker/uuids.txt').split('\n').filter(u => u.trim().length > 0);
});

export default function () {
    const uuid = uuids[Math.floor(Math.random() * uuids.length)];
    const res = http.get(`http://localhost:8888/items/${uuid}`);
}