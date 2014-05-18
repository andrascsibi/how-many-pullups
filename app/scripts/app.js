angular.module('pullApp', [
  'ngRoute',
  'ngResource',
  'mgcrea.ngStrap',
])

.factory('ValidatorService', ['$q', '$http', '$route', 'baseUrl', function($q, $http, $route, baseUrl) {
  return function() {
    var delay = $q.defer();
    var promise = ($route.current.params.jsonFile) ?
      $http.get($route.current.params.jsonFile) :
      $http.jsonp(baseUrl + '?callback=JSON_CALLBACK');
    promise.success(function(data, status, headers, config) {
      delay.resolve(data);
    }).
    error(function(data, status, headers, config) {
      delay.reject('Unable to fetch validator data');
    });
    return delay.promise;
  };
}])


.config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
  'use strict';
  $locationProvider.html5Mode(true);

  $routeProvider
    .when('/', {
      controller:'HomePageCtrl',
      templateUrl:'/app/views/index.html'
    })
    .when('/andris', {
      controller:'TotalCtrl',
      templateUrl:'/app/views/jumbocounter.html'
    })
    .when('/admin/accounts', {
      controller:'AdminCtrl',
      templateUrl:'/app/views/admin/accounts.html'
    })
    .when('/:id', {
      controller:'BoardCtrl',
      templateUrl:'/app/views/board.html'
    })
    .otherwise({
      redirectTo:'/'
    });
}]);
