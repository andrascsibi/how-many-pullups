console.log("what up");

angular.module('how-many-pullups', ['ngRoute'])
 
.value('reloadInterval', 5 * 60 * 1000)
 
//.factory('Projects', function($firebase, fbURL) {
//  return $firebase(new Firebase(fbURL));
//})
 
.config(['$routeProvider', '$locationProvider',
  function($routeProvider, $locationProvider) {
  $locationProvider.html5Mode(true);

  $routeProvider
    .when('/', {
      controller:'TotalCtrl',
      templateUrl:'/html/jumbocounter.html'
    })
    .when('/admin', {
       controller:'TotalCtrl',
       templateUrl:'/html/admin.html'
    })
    .when('/hello', {
       controller:'HelloCtrl',
       templateUrl:'/html/hello.html'
    })
    .otherwise({
      redirectTo:'/'
    });
}])
 
.controller('TotalCtrl', ['$scope', '$http', '$interval', 'reloadInterval',
  function($scope, $http, $interval, reloadInterval) {

  if (!$scope.refresh) {
    var refresh = function() {
      $http({method: 'GET', url: 'total'}).
        success(function(data, status, headers, config) {
          $scope.stat = data;
        }).
        error(function(data, status, headers, config) {
          console.log("request failed");
      });
    };
    $scope.refresh = refresh;
    refresh();
    $interval(refresh, reloadInterval);
  }
}])
 
.controller('HelloCtrl', ['$scope', '$http', function($scope, $http) {
  $http({method: 'GET', url: 'whoami'}).
    success(function(data, status, headers, config) {
      $scope.stat = data;
    }).
    error(function(data, status, headers, config) {
      console.log("request failed");
  });
}]);
 
 
