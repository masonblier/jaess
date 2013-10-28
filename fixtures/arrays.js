var a = [1, null, 3];

for (var i = 0; i < 3; ++i) {
  if (a[i] !== undefined) {
    delete a;
  }
}

a[2];
