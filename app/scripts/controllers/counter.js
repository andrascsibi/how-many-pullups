angular.module('pullApp')

.controller('CounterCtrl', ['$scope', '$http', 'Stats',
  function($scope, $http, Stats){
  var accountId = $scope.c.AccountID;
  var challengeId = $scope.c.ID;
  var statUrl = '/accounts/' + accountId + '/challenges/' + challengeId + '/stats';

  Stats.get({id: accountId, c_id: challengeId}, function(stat){
    $scope.stat = stat;
  }, function(error){
    alert(error.data.error); // TODO
  });


  $http.get(statUrl).success(function(stat) {
  });
  $scope.foo = "bar";

}]);
