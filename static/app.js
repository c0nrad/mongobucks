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
      controller: 'CashoutController',
      templateUrl: "partials/cashout.html",
    })
    .state('redeem', {
      url: "/r/:token",
      controller: 'RedemptionController',
      templateUrl: 'partials/redemption.html'
    })
    .state('me', {
      url: "/me",
      controller: 'MeController',
      templateUrl: 'partials/me.html'
    })
});

app.service("User", function($resource) {
  return $resource("/api/users/:user", {}, {
    me: { method: "get", "url": "/api/users/me"},
    tickets: {method: "get", isArray:true, url: "/api/users/me/tickets" }
  })
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

app.service("Reward", function($resource) {
  return $resource("/api/rewards/:id", {"id": "@_id"}, { 
    recent: { method: "get", isArray:true,  url: "/api/gambles/recent"},
    user: {method: "get", isArray: true, url: "/api/users/:user/gambles" }
  })
});

app.service("Ticket", function($resource) {
  return $resource("/api/tickets/:token", {}, {
    redeem: {method: "POST", url: "/api/tickets/:token/redeem"}
  });
});


app.controller("HomeController", function($scope, User, Transaction, Gamble) {
  $scope.users = User.query();
  $scope.recentTransactions = Transaction.recent();
  $scope.recentGambles = Gamble.recent()
})

app.controller("CashoutController", function($scope, $stateParams, $state, Reward, Ticket, User) {
  $scope.rewards = Reward.query()
  $scope.me = User.me();


  $scope.buy = function(reward) {
    Ticket.save({reward: reward.ID}, function() {
      $state.go('me');
    });
  }
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

app.controller("RedemptionController", function($scope, $state, $stateParams, User, Ticket) {
  $scope.me = User.me();
  $scope.ticket = Ticket.get({token: $stateParams.token})

  $scope.redeem = function() {
    $scope.ticket.$redeem({token: $scope.ticket.Redemption}, function() {
      $state.reload();
    })
  }
})


app.controller("UserController", function($scope, $stateParams, User, Transaction, Gamble) {
  $scope.user = User.get({user: $stateParams.user})
  $scope.transactions = Transaction.user({user: $stateParams.user})
  $scope.gambles = Gamble.user({user: $stateParams.user})
})

app.controller('MeController', function($scope, User) {
  $scope.me = User.me()
  $scope.tickets = User.tickets()
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

