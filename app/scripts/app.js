angular.module('pullApp', [
  'ngRoute'
])
.config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
  'use strict';
  $locationProvider.html5Mode(true);

  $routeProvider
    .when('/', {
      controller:'TotalCtrl',
      templateUrl:'/app/views/jumbocounter.html'
    })
    .when('/admin', {
       controller:'TotalCtrl',
       templateUrl:'/app/views/admin.html'
    })
    .when('/hello', {
       controller:'HelloCtrl',
       templateUrl:'/app/views/hello.html'
    })
    .otherwise({
      redirectTo:'/'
    });
}]);
