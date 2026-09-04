import{Fi as e,Jr as t,Kn as n,Li as r,Mn as i,Oi as a,Qn as o,Yr as s,kn as c,qn as l}from"./router-Ca5XLGLk.js";import{n as u}from"./index-ChrUCjfM.js";var d=t(`statistic`,[s(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),t(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[s(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[t(`icon`,{verticalAlign:`-0.125em`})]),s(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),s(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[t(`icon`,{verticalAlign:`-0.125em`})])])]),f=e({name:`Statistic`,props:Object.assign(Object.assign({},c.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:r,mergedRtlRef:o}=l(e),s=c(`Statistic`,`-statistic`,d,u,e,t),f=i(`Statistic`,o,t),p=a(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:c},common:{cubicBezierEaseInOut:l}}=s.value;return{"--n-bezier":l,"--n-label-font-size":c,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=r?n(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:t,cssVars:r?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:n,label:i,prefix:a,suffix:s}}=this;return(e=this.onRender)==null||e.call(this),r(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},o(i,e=>r(`div`,{class:`${t}-statistic__label`},this.label||e)),r(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},o(a,e=>e&&r(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?o(n,e=>e&&r(`span`,{class:`${t}-statistic-value__content`},e)):r(`span`,{class:`${t}-statistic-value__content`},this.value),o(s,e=>e&&r(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};