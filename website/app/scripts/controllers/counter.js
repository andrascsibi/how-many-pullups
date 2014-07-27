angular.module('pullApp')

.controller('CounterCtrl', ['$scope', '$http', 'Stats',
  function($scope, $http, Stats){
  var accountId = $scope.challenge.AccountID;
  var challengeId = $scope.challenge.ID;

  var defaultDateToNull = function(date) {
    return date === '0001-01-01T00:00:00Z' ? null : date;
  };

  $scope.get = function() {
    Stats.get({id: accountId, c_id: challengeId}, function(stat){
      stat.MinDate = defaultDateToNull(stat.MinDate);
      stat.MaxDate = defaultDateToNull(stat.MaxDate);
      $scope.stat = stat;
    }, function(error){
      alert(error.data.error); // TODO
    });
  };

  $scope.get();

}]);
