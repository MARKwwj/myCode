class Test_case:
    def test_one(self):
        x = "this"
        assert "h" in x

    def test_two(self):
        f = "helloworld"
        assert f == "helloworlds"


def test_passing():
    assert (1, 2, 3) == (1, 2, 3)
