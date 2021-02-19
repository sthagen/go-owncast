var request = require('supertest');
request = request('http://127.0.0.1:8080');

var accessToken = '';
var webhookID;

const webhook = 'https://super.duper.cool.thing.biz/owncast';
const events = ['CHAT'];

test('create webhook', async (done) => {
    const res = await sendIntegrationsChangePayload('webhooks/create', {
        url: webhook,
        events: events, 
    });

    expect(res.body.url).toBe(webhook);
    expect(res.body.timestamp).toBeTruthy();
    expect(res.body.events).toStrictEqual(events);
    done();
});

test('check webhooks', (done) => {
    request.get('/api/admin/webhooks')
        .auth('admin', 'abc123').expect(200)
        .then((res) => {
            expect(res.body).toHaveLength(1);
            expect(res.body[0].url).toBe(webhook);
            expect(res.body[0].events).toStrictEqual(events);
            webhookID = res.body[0].id;
            done();
        });
});

test('delete webhook', async (done) => {
    const res = await sendIntegrationsChangePayload('webhooks/delete', {
        id: webhookID,
    });
    expect(res.body.success).toBe(true);
    done();
});

test('check that webhook was deleted', (done) => {
    request.get('/api/admin/webhooks')
        .auth('admin', 'abc123').expect(200)
        .then((res) => {
            expect(res.body).toHaveLength(0);
            done();
        });
});

test('create access token', async (done) => {
    const name = 'test token';
    const scopes = ['CAN_SEND_SYSTEM_MESSAGES'];
    const res = await sendIntegrationsChangePayload('accesstokens/create', {
        name: name,
        scopes: scopes, 
    });

    expect(res.body.token).toBeTruthy();
    expect(res.body.timestamp).toBeTruthy();
    expect(res.body.name).toBe(name);
    expect(res.body.scopes).toStrictEqual(scopes);
    accessToken = res.body.token;
    done();
});

test('check access tokens', (done) => {
    request.get('/api/admin/accesstokens')
        .auth('admin', 'abc123').expect(200)
        .then((res) => {
            expect(res.body).toHaveLength(1);
            expect(res.body[0].token).toBe(accessToken);
            done();
        });
});

test('send a system message using access token', async (done) => {
    const payload = {body: 'test 1234'};
    const res = await request.post('/api/integrations/chat/system')
        .set('Authorization', 'Bearer ' + accessToken)
        .send(payload).expect(200);
    done();
});

test('delete access token', async (done) => {
    const res = await sendIntegrationsChangePayload('accesstokens/delete', {
        token: accessToken,
    });
    expect(res.body.success).toBe(true);
    done();
});

test('check token delete was successful', (done) => {
    request.get('/api/admin/accesstokens')
        .auth('admin', 'abc123').expect(200)
        .then((res) => {
            expect(res.body).toHaveLength(0);
            done();
        });
});

async function sendIntegrationsChangePayload(endpoint, payload) {
    const url = '/api/admin/' + endpoint;
    const res = await request.post(url)
        .auth('admin', 'abc123')
        .send(payload).expect(200);

    return res
}