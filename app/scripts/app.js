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
      title: 'Counting Made Easy',
      controller: 'HomepageCtrl',
      templateUrl:'/app/views/index.html',
    })
    .when('/admin/accounts', {
      title: 'Accounts',
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
