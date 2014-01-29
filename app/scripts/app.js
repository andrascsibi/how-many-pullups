angular.module('pullApp', [
  'ngRoute'
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
