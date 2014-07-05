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
      resolve: {
        whoami: ['WhoamiService', function(WhoamiService) {
          return WhoamiService();
        }],
      },
    })
    .when('/admin/accounts', {
      title: 'Accounts',
      controller:'AdminCtrl',
      templateUrl:'/app/views/admin/accounts.html',
    })
    .when('/:id', {
      controller:'BoardCtrl',
      templateUrl:'/app/views/board.html',
      resolve: {
        whoami: ['WhoamiService', function(WhoamiService) {
          return WhoamiService();
        }],
      },
    })
    .when('/:id/:c_id', {
      controller:'ChallengeCtrl',
      templateUrl:'/app/views/challenge.html',
      resolve: {
        allSets: ['AllSets', '$route', function(AllSets, $route) {
          return AllSets.query($route.current.params).$promise;
        }],
        challenge: ['Challenge', '$route', function(Challenge, $route) {
          return Challenge.get($route.current.params).$promise;
        }],
        whoami: ['WhoamiService', function(WhoamiService) {
          return WhoamiService();
        }],
      },
    })
    .otherwise({
      redirectTo:'/'
    });
}]);
