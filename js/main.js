console.log("what up");

angular.module('how-many-pullups', ['ngRoute'])
 
//.value('fbURL', 'https://angularjs-projects.firebaseio.com/')
 
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
 
.controller('TotalCtrl', ['$scope', '$http', function($scope, $http) {
  
  $http({method: 'GET', url: 'total'}).
    success(function(data, status, headers, config) {
      $scope.stat = data;
    }).
    error(function(data, status, headers, config) {
      console.log("request failed");
  });
}])
 
.controller('HelloCtrl', ['$scope', '$http', function($scope, $http) {
  $scope.email = "foo";
}]);
 
 
