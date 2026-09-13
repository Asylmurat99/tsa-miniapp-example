const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

// #region check
// fields: decoded values; signField: "hash" for the launch context, "sign"
// for the getPhone envelope.
function checkSignature(fields, secret, signField) {
  const { [signField]: received = '', ...rest } = fields;
  const canon = Object.keys(rest).sort().map((k) => `${k}=${rest[k]}`).join('\n');
  const calc = crypto.createHmac('sha256', secret).update(canon, 'utf8').digest('hex');
  return calc.length === received.length &&
    crypto.timingSafeEqual(Buffer.from(calc), Buffer.from(received));
}
// #endregion check

const vector = JSON.parse(fs.readFileSync(path.join(__dirname, '..', 'vector.json'), 'utf8'));
let failed = false;
for (const c of vector.cases) {
  const ok = checkSignature(c.fields, vector.secret, c.sign_field);
  console.log(`${c.name}: ${ok ? 'ok' : 'FAIL'}`);
  if (!ok) failed = true;
}
process.exit(failed ? 1 : 0);
