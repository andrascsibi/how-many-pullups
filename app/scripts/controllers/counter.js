angular.module('pullApp')

.controller('CounterCtrl', ['$scope', '$http', 'Stats',
  function($scope, $http, Stats){
  var accountId = $scope.challenge.AccountID;
  var challengeId = $scope.challenge.ID;

  $scope.get = function() {
    Stats.get({id: accountId, c_id: challengeId}, function(stat){
      $scope.stat = stat;
    }, function(error){
      alert(error.data.error); // TODO
    });
  };

  $scope.get();

}]);
