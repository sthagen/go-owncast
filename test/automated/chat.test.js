const { test } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const WebSocket = require('ws');
var ws;

const testMessageId = Math.random().toString(36).substring(7);
const username = 'user' + Math.floor(Math.random() * 100);
const message = Math.floor(Math.random() * 100) + ' test 123';
const messageRaw = message + ' *and some markdown too*';
const messageMarkdown = '<p>' + message + ' <em>and some markdown too</em></p>'
const date = new Date().toISOString();

const testMessage = {
    author: username,
    body: messageRaw,
    id: testMessageId,
    type: 'CHAT',
    visible: true,
    timestamp: date,
};

test('can send a chat message', (done) => {
    ws = new WebSocket('ws://127.0.0.1:8080/entry', {
    origin: 'http://localhost',
});

    function onOpen() {
        ws.send(JSON.stringify(testMessage), function() {
            ws.close();
            done();
        });
    }

    ws.on('open', onOpen);
});

test('can fetch chat messages', (done) => {
    request.get('/api/admin/chat/messages').auth('admin', 'abc123').expect(200)
        .then((res) => {
            const message = res.body.filter(function(msg) {
                return msg.id = testMessageId;
            })[0];

            expect(message.author).toBe(testMessage.author);
            expect(message.body).toBe(messageMarkdown);
            expect(message.date).toBe(testMessage.date);
            expect(message.type).toBe(testMessage.type);

            done();
        });
});
