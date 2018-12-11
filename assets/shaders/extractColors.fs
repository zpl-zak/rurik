#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;

uniform vec3 treshold = vec3(0.4);

// Output fragment color
out vec4 finalColor;

void main()
{
    vec4 col = texture2D(texture0, fragTexCoord);

    float lum = dot(col.rgb, treshold);

    if (lum > 1.0) {
        finalColor = col;
    } else {
        finalColor = vec4(0.0, 0.0, 0.0, 1.0);
    }
}