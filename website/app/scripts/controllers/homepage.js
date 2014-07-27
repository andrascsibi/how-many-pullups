angular.module('pullApp')

.controller('HomepageCtrl', ['$scope', 'whoami', '$routeParams', '$location', '$modal', 'Account', '$timeout',
  function($scope, whoami, $routeParams, $location, $modal, Account, $timeout){

  $scope.whoami = whoami;

  var regModal = $modal({scope: $scope, template: 'app/views/registration.html', show: false});
  $scope.showRegModal = function() {
    regModal.$promise.then(regModal.show);
  };

  $scope.createAccount = function(account) {
    var newAccount = new Account();
    newAccount.Email = account.Email;
    newAccount.ID = account.ID;
    $scope.working = true;
    newAccount.$save(function(a, putRespHeaders) {
      $timeout(function() {
        $scope.working = false;
        $scope.error = null;
        $scope.whoami.Account = a;
        regModal.hide();
        $location.path('/' + a.ID);
      }, 2000); // XXX: dirtiest trick in the book
    }, function(err) {
      $scope.working = false;
      $scope.error = err.data;
    });
  };

  if ($scope.whoami.Unregistered) {
    console.log("showing reg modal");
    $scope.showRegModal();
  }
  if ($routeParams.redir && !whoami.Unregistered) {
    $location.path('/' + whoami.Account.ID).search({});
  }


}]);
