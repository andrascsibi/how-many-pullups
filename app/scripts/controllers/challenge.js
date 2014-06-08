angular.module('pullApp')

.controller('ChallengeCtrl', ['$scope', '$rootScope', 'challenge', 'allSets', 'WhoamiService',
  function($scope, $rootScope, challenge, allSets, WhoamiService) {

  $scope.challenge = challenge;
  $rootScope.title = challenge.AccountID + "'s " + challenge.Title + ' challenge';

  var parseDates = function(sets) {
    return sets
    .filter(function(cur) {
      return cur !== 0;
    })
    .map(function(cur) {
      return {
        ts: Math.round(Date.parse(cur.Date)/1000),
        reps: cur.Reps
      };
    });
  };

  var toCalHeatmap = function(sets) {
    return sets.reduce(function(memo, cur) {
      if (!memo[cur.ts]) {
        memo[cur.ts] = 0;
      }
      memo[cur.ts] += cur.reps;
      return memo;
    }, {});
  };

  var getStats = function(sets) {
    if (sets.length === 0) {
      return null;
    }
    return sets.reduce(function(memo, cur) {
      memo.numSets++;
      memo.totalReps += cur.Reps;
      memo.maxReps = Math.max(memo.maxReps, cur.Reps);
      return memo;
    }, {
      numSets: 0,
      totalReps: 0,
      maxReps: 0,
      minDate: sets[sets.length - 1].Date,
      maxDate: sets[0].Date,
    });
  };

  var groupByDay = function(sets) {
    return sets.reduce(function(memo, cur) {
      var date = new Date(cur.ts*1000).toDateString();
      if (!memo[date]) {
        memo[date] = {sets: [], reps: 0};
      }
      memo[date].sets.push(cur);
      memo[date].reps += cur.reps;
      return memo;
    }, {});
  };

  var sets = parseDates(allSets);
  $scope.stats = getStats(allSets);

  var setsByDay = groupByDay(sets);
  var repsByDay = Object.keys(setsByDay).map(function(k) {
    return setsByDay[k].reps;
  });
  repsByDay.sort(function(l,r) {return l-r;});

  var getPercentile = function(arr, p) {
    arr.sort(function(l,r) {return l-r;});
    return arr[Math.floor(arr.length * p)];
  };

  var getLegend = function(max) {
    var step = max / 4;
    return [1, 2, 3, 4].map(function(e){
      return Math.round(e * step);
    });
  };

  var cal = new CalHeatMap();
  cal.init({
    start: new Date($scope.stats.minDate),
    range: 6,
    end: new Date(),
    itemName: ['rep', 'reps'],
    tooltip: true,
    domain: "month",
    subDomain: "day",
    data: toCalHeatmap(sets),
    cellSize: 15,
    legendCellSize: 15,
    legend: getLegend(getPercentile(repsByDay, 0.9)),
  });

}]);
