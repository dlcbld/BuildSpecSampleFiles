from flask import Flask
from flask import request

app = Flask(__name__)

@app.route("/")
def index():
    pound = request.args.get("pound", "")
    if pound:
        kilogram = kilogram_from(pound)
    else:
        kilogram = ""
    return (
        """<form action="" method="get">
                Weight in pounds: <input type="text" name="pound">
                <input type="submit" value="Convert to Kilograms">
            </form>"""
        + "Weight in Kilograms: "
        + kilogram)

def kilogram_from(pound):
    """Convert Weight in Pounds to Weight in Kilograms."""
    try:
        kilogram = float(pound) * 0.454
        kilogram = round(kilogram, 3) 
        return str(kilogram)
    except ValueError:
        return "invalid input"

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=8080, debug=True)
