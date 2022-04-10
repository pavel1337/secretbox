import http from 'k6/http';
import { parseHTML } from 'k6/html';
import { sleep } from 'k6';

const urlCreateSecret = 'http://localhost:4000/secret/create';
const baseUrl = 'http://localhost:4000';

export const options = {
    scenarios: {
      constant_request_rate: {
        executor: 'constant-arrival-rate',
        rate: 2000,
        timeUnit: '1s', // 1000 iterations per second, i.e. 1000 RPS
        duration: '30s',
        preAllocatedVUs: 100, // how large the initial pool of VUs would be
        maxVUs: 500, // if the preAllocatedVUs are not enough, we can initialize more
      },
    },
  };

export default function () {
  let data = { content: 'secret', expires: "10" };


  // Using an object as body, the headers will automatically include
  // 'Content-Type: application/x-www-form-urlencoded'.
  let res = http.post(urlCreateSecret, data);
  let doc = parseHTML(res.body);
  let link = doc.find('div.flash a').attr('href')

//   console.log(link); // Bert

  http.get(baseUrl+link)

  data = { continue: 'true' };
  http.post(baseUrl+link, data);
}
