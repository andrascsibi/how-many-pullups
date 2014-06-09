angular.module('pullApp')

.controller('ChallengeCtrl', ['$scope', '$rootScope', 'challenge', 'allSets', 'whoami',
  function($scope, $rootScope, challenge, allSets, whoami) {

  $scope.challenge = challenge;
  $rootScope.title = challenge.AccountID + "'s " + challenge.Title + ' challenge';
  $scope.whoami = whoami;

  if (allSets.length === 0) {
    $scope.empty = true;
    return;
  }

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

  $scope.stats = getStats(allSets);

  var sets = parseDates(allSets);

  var dayKey = function(set) {
    return new Date(set.ts*1000)
        .toISOString()
        .substring(0, "2014-06-09".length);
  };
  var hourKey = function(set) {
    return new Date(set.ts*1000)
        .toISOString()
        .substring(0, "2014-06-09T20".length);
  };

  var histogram = function(sets, keyFun) {
    var groups = _.groupBy(sets, keyFun);
    var repCount = _.values(groups)
        .map(function(g) {
          return g.reduce(function(sum, cur) {
            return sum + cur.reps;
          }, 0);
        });
    return _.sortBy(repCount, function(l,r) {return l-r;});
  };

  var repsByDay = histogram(sets, dayKey);
  var repsByHour = histogram(sets, hourKey);

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

  $scope.stats.workDays = repsByDay.length;

  var calSettings = {
    start: new Date($scope.stats.minDate),
    end: new Date(),
    range: 6,
    domain: "week",
    itemName: ['rep', 'reps'],
    tooltip: true,
    data: toCalHeatmap(sets),
    cellSize: 15,
    legendCellSize: 15,
    domainGutter: 10,
  };

  var cal = new CalHeatMap();
  cal.init(angular.extend({
    itemSelector: '#cal-heatmap-punchcard',
    itemNamespace: 'punchcard',
    rowLimit: 24,
    subDomain: "hour",
    legend: getLegend(getPercentile(repsByHour, 0.9)),
    legendHorizontalPosition: 'right',
    legendVerticalPosition: 'top',
    label: {position: 'top'},
  }, calSettings));

  var cal2 = new CalHeatMap();
  cal2.init(angular.extend({
    itemSelector: '#cal-heatmap-daily',
    itemNamespace: 'daily',
    rowLimit: 1,
    subDomain: "day",
    subDomainTextFormat: function(date, value) {
      return value;
    },
    legend: getLegend(getPercentile(repsByDay, 0.9)),
    legendHorizontalPosition: 'right',
    legendVerticalPosition: 'bottom',
    domainLabelFormat: '',
  }, calSettings));

  $scope.next = function() {
    cal.next();
    cal2.next();
  };

  $scope.previous = function() {
    cal.previous();
    cal2.previous();
  };
}]);
