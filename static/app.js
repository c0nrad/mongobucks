var app = angular.module("app", ["ui.router", "ngResource"]);

app.config(function($stateProvider, $urlRouterProvider) {
  //
  // For any unmatched url, redirect to /state1
  $urlRouterProvider.otherwise("/");
  //
  // Now set up the states
  $stateProvider
    .state('home', {
      url: "/",
      templateUrl: "partials/home.html",
      controller: "HomeController"
    })
    .state('transaction', {
      url: "/t/:id",
      templateUrl: "partials/transaction.html",
      controller: "TransactionController"
    })
    .state('gamble', {
      url: "/g/:id",
      templateUrl: "partials/gamble.html",
      controller: "GambleController"
    })
    .state('user', {
      url: "/u/:user",
      templateUrl: "partials/user.html",
      controller: "UserController"
    })
    .state('cashout', {
      url: "/cashout",
      templateUrl: "partials/cashout.html",
    })
});

app.service("User", function($resource) {
  return $resource("/api/users/:user", {}, {me: { method: "get", "url": "/api/users/me"}})
});

app.service("Transaction", function($resource) {
  return $resource("/api/transactions/:id", {"id": "@_id"}, { 
    recent: { method: "get", isArray:true,  url: "/api/transactions/recent"},
    user: {method: "get", isArray: true, url: "/api/users/:user/transactions" }
  })
});

app.service("Gamble", function($resource) {
  return $resource("/api/gambles/:id", {"id": "@_id"}, { 
    recent: { method: "get", isArray:true,  url: "/api/gambles/recent"},
    user: {method: "get", isArray: true, url: "/api/users/:user/gambles" }
  })
});

app.controller("HomeController", function($scope, User, Transaction, Gamble) {
  $scope.users = User.query();
  $scope.recentTransactions = Transaction.recent();
  $scope.recentGambles = Gamble.recent()
})

app.controller("TransactionController", function($scope, $stateParams, Transaction) {
  $scope.transaction = Transaction.get({id: $stateParams.id})
})

app.controller("GambleController", function($scope, $stateParams, Gamble) {
  $scope.gamble = Gamble.get({id: $stateParams.id})
})


app.controller("HeaderController", function($scope, User) {
  $scope.me = User.me();
})

app.controller("UserController", function($scope, $stateParams, User, Transaction, Gamble) {
  $scope.user = User.get({user: $stateParams.user})
  $scope.transactions = Transaction.user({user: $stateParams.user})
  $scope.gambles = Gamble.user({user: $stateParams.user})
})


//<!-- Filters -->
app.filter('fromNow', function() {
  return function(date) {
    if (moment(date).isBefore(moment("2000-01-01T00:00:00Z")))  {
      return "";
    } 
    return moment(date).fromNow();
  };
});

app.filter('duration', function() {
  return function(date) {
    return moment.duration(moment().diff(date)).humanize()
  }
})

