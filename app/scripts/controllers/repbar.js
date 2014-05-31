angular.module('pullApp')

.controller('RepbarCtrl', ['$scope', 'Sets', function($scope, Sets) {

  var range = function(from, to, step) {
    var range = [];
    for (var i = from; i <= to; i += step) {
      range.push(i);
    }
    return range;
  };

  $scope.repButtons = range(1, $scope.challenge.MaxReps || 10, 1);

  $scope.add = function(reps) {

    $scope.working = true;
    var newSet = new Sets();
    newSet.Reps = reps;
    newSet.$save({id: $scope.challenge.AccountID, c_id: $scope.challenge.ID})
    .then(function(){
      $scope.working = false;
      $scope.list();
    }, function(err) {
      alert(err.data.error); // TODO
    });
  };

}]);

