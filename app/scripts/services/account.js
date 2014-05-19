angular.module('pullApp')

.factory("Account", ['$resource', function ($resource) {
  return $resource("/accounts/:id", {id: "@ID" }, {});
}]);
