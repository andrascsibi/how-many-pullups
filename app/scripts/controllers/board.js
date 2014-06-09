angular.module('pullApp')

.controller('BoardCtrl', ['$scope', '$rootScope', 'Account', 'Challenge', '$routeParams', 'whoami',
  function($scope, $rootScope, Account, Challenge, $routeParams, whoami) {

  $rootScope.title = $routeParams.id + "'s challenges";

  $scope.whoami = whoami;
  $scope.owner = whoami.owner;

  Account.get({id: $routeParams.id}, function(data){
    $scope.account = data;
  }, function(err) {
    if (err.status === 404) {
      $scope.notFound = true;
    }
  });

  $scope.challenges = [];

  $scope.list = function() {
    Challenge.query({id: $routeParams.id}, function(data){
      $scope.challenges = data;
    }, function(error){
      alert(error.data.error); // TODO
    });
  };

  $scope.list();

  $scope.edited = null;

  $scope.add = function() {
    var newChallenge = new Challenge();
    newChallenge.Title = 'Pullups';
    newChallenge.Description = '';
    newChallenge.MaxReps = 10;
    newChallenge.StepReps = 1;
    newChallenge.AccountID = $routeParams.id;
    $scope.challenges.splice(0, 0, newChallenge);
    $scope.edited = newChallenge;
  };

  $scope.edit = function(c) {
    $scope.edited = angular.copy(c);
  };

  $scope.cancel = function() {
    $scope.edited = null;
    $scope.list();
  };

  $scope.save = function() {
    $scope.working = true;
    $scope.edited.$save()
    .then(function(){
    }, function(err) {
      alert(err.data.error);
    })
    .then(function(){
      $scope.working = false;
      $scope.edited = null;
      $scope.list();
    });
  };


}]);

