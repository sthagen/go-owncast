module.exports = { createTestMessageObject };

function createTestMessageObject(userContext, events, done) {
  const randomNumber = Math.floor((Math.random() * 10) + 1);
  const author = "load-test-user-" + randomNumber
  const data = {
    author: author,
    body: "Test 12345. " + randomNumber,
    image: "https://robohash.org/" + author + "?size=80x80&set=set3"
  };
  // set the "data" variable for the virtual user to use in the subsequent action
  userContext.vars.data = data;
  return done();
}
