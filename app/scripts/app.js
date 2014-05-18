angular.module('pullApp', [
  'ngRoute',
  'ngResource',
  'mgcrea.ngStrap',
])

.factory('WhoamiService', ['$q', '$http', '$route', function($q, $http, $route) {
  return function() {
    var delay = $q.defer();
    var promise = $http.get('/whoami');
    promise.success(function(data, status, headers, config) {
      delay.resolve(data);
    }).
    error(function(data, status, headers, config) {
      delay.reject('Could not fetch whoami');
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
      templateUrl:'/app/views/index.html',
    })
    .when('/admin/accounts', {
      controller:'AdminCtrl',
      templateUrl:'/app/views/admin/accounts.html',
    })
    .when('/:id', {
      controller:'BoardCtrl',
      templateUrl:'/app/views/board.html',
    })
    .otherwise({
      redirectTo:'/'
    });
}]);
