angular.module('pullApp', [
  'ngRoute',
  'ngResource',
  'mgcrea.ngStrap',
])

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
