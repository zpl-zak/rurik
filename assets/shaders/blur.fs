#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;

// Output fragment color
out vec4 finalColor;

uniform bool horizontal;
uniform vec2 size = vec2(640, 480);

float weight[3] = float[](0.2270270270, 0.3162162162, 0.0702702703);

void main()
{
    // Texel color fetching from texture sampler
    vec3 texelColor = texture(texture0, fragTexCoord).rgb*weight[0];
    vec2 texelOffset = 1.0 / textureSize(texture0, 0);

    if (horizontal) {
        for (int i = 1; i < 5; ++i) {
            texelColor += texture(texture0, fragTexCoord + vec2(texelOffset.x * i, 0.0)).rgb * weight[i];
            texelColor += texture(texture0, fragTexCoord - vec2(texelOffset.x * i, 0.0)).rgb * weight[i];
        }
    } else {
        for (int i = 1; i < 5; ++i) {
            texelColor += texture(texture0, fragTexCoord + vec2(0.0, texelOffset.x * i)).rgb * weight[i];
            texelColor += texture(texture0, fragTexCoord - vec2(0.0, texelOffset.x * i)).rgb * weight[i];
        }
    }

    finalColor = vec4(texelColor, 1.0);
}