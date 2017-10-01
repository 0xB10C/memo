var express = require('express');
var path = require('path');
var logger = require('morgan');

var index = require('./routes/index');
var api = require('./routes/api');


var app = express();

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');

app.use(logger('combined'));
app.use(express.static(path.join(__dirname, 'public')));

app.use('/', index);
app.use('/api', api);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  var err = new Error('Not Found');
  err.status = 404;
  next(err);
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  //res.locals.error = req.app.get('env') === 'development' ? err : {};

  console.log(res.locals.message);
  console.log(res.locals.err);

  // render the error page
  //res.status(err.status || 500);
  res.sendFile(__dirname + "/views/error.html");
});

module.exports = app;
